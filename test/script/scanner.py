import tkinter as tk
from tkinter import ttk, messagebox, scrolledtext
import requests
import threading
import ipaddress

class GatewayScannerApp:
    def __init__(self, root):
        self.root = root
        self.root.title("局域网网关扫描器")

        # 扫描网段选择框
        self.network_label = tk.Label(root, text="选择扫描网段:")
        self.network_label.pack(pady=10)
        self.network_var = tk.StringVar()
        self.network_combobox = ttk.Combobox(root, textvariable=self.network_var,
                                             values=["192.168.1.0/24", "192.168.0.0/24", "10.0.0.0/24"])
        self.network_combobox.current(0)
        self.network_combobox.pack(pady=5)

        # 扫描按钮
        self.scan_button = tk.Button(root, text="开始扫描", command=self.start_scan)
        self.scan_button.pack(pady=20)

        # 停止扫描按钮
        self.stop_button = tk.Button(root, text="停止扫描", command=self.stop_scan, state=tk.DISABLED)
        self.stop_button.pack(pady=10)

        # 列表框及滚动条
        self.listbox = tk.Listbox(root, width=80, height=20)
        self.listbox.pack(side=tk.LEFT, fill=tk.BOTH, expand=True)
        self.scrollbar = tk.Scrollbar(root)
        self.scrollbar.pack(side=tk.RIGHT, fill=tk.Y)
        self.listbox.config(yscrollcommand=self.scrollbar.set)
        self.scrollbar.config(command=self.listbox.yview)

        # 鼠标经过列表项变色
        self.listbox.bind("<Enter>", self.on_enter)
        self.listbox.bind("<Leave>", self.on_leave)

        # 日志输出框
        self.log_text = scrolledtext.ScrolledText(root, width=80, height=10)
        self.log_text.pack(pady=10)

        # 整体进度条
        self.progress_label = tk.Label(root, text="扫描进度: 0%")
        self.progress_label.pack(pady=5)

        # 存储网关信息
        self.gateway_info = {}

        # 绑定列表项点击事件
        self.listbox.bind("<Double-1>", self.on_listbox_click)

        # 扫描线程列表
        self.threads = []
        # 停止标志
        self.stop_event = threading.Event()

    def start_scan(self):
        network_str = self.network_var.get()
        try:
            # 验证输入的网段格式
            network = ipaddress.ip_network(network_str, strict=False)
        except ValueError:
            messagebox.showerror("错误", "无效的网段格式，请输入正确的网段，例如：192.168.1.0/24")
            return

        # 清空列表框、网关信息和日志
        self.listbox.delete(0, tk.END)
        self.gateway_info = {}
        self.log_text.delete(1.0, tk.END)
        self.progress_label.config(text="扫描进度: 0%")

        # 重置停止标志
        self.stop_event.clear()

        # 启用停止按钮
        self.scan_button.config(state=tk.DISABLED)
        self.stop_button.config(state=tk.NORMAL)

        # 获取网段内所有 IP 地址
        ip_list = list(network.hosts())
        total_ips = len(ip_list)

        # 启动多线程扫描
        num_threads = 10  # 线程数量
        chunk_size = total_ips // num_threads
        for i in range(num_threads):
            start = i * chunk_size
            end = start + chunk_size if i < num_threads - 1 else total_ips
            thread = threading.Thread(target=self.scan_range, args=(ip_list[start:end], total_ips))
            self.threads.append(thread)
            thread.start()

    def scan_range(self, ip_range, total_ips):
        scanned_count = 0
        for index, ip in enumerate(ip_range):
            if self.stop_event.is_set():
                break
            scanned_count += 1
            progress = int((scanned_count / total_ips) * 100)
            self.log_text.insert(tk.END, f"正在扫描: {ip} ({scanned_count}/{total_ips})\n")
            self.log_text.see(tk.END)  # 自动滚动到最新日志
            self.progress_label.config(text=f"扫描进度: {progress}%")
            url = f"http://{ip}:2580/api/v1/os/system"
            try:
                response = requests.get(url, timeout=2)
                if response.status_code == 200:
                    data = response.json()
                    if data.get("code") == 200:
                        info = f"IP: {ip}, CPU 使用率: {data['data']['hardWareInfo']['cpuPercent']}%, 磁盘使用率: {data['data']['hardWareInfo']['diskInfo']}%, 内存使用率: {data['data']['hardWareInfo']['memPercent']}%"
                        self.listbox.insert(tk.END, info)
                        self.gateway_info[len(self.gateway_info)] = str(ip)
                        self.log_text.insert(tk.END, f"发现网关: {ip}\n")
                        self.log_text.see(tk.END)  # 自动滚动到最新日志
            except requests.RequestException:
                continue

        if not self.stop_event.is_set():
            self.log_text.insert(tk.END, "扫描完成\n")
            self.log_text.see(tk.END)  # 自动滚动到最新日志
            self.scan_button.config(state=tk.NORMAL)
            self.stop_button.config(state=tk.DISABLED)

    def stop_scan(self):
        self.stop_event.set()
        self.scan_button.config(state=tk.NORMAL)
        self.stop_button.config(state=tk.DISABLED)

    def on_listbox_click(self, event):
        index = self.listbox.curselection()
        if index:
            ip = self.gateway_info[index[0]]
            import webbrowser
            webbrowser.open(f"http://{ip}:2580")

    def on_enter(self, event):
        # 鼠标进入列表项时变色
        index = self.listbox.nearest(event.y)
        self.listbox.itemconfig(index, bg='lightblue')

    def on_leave(self, event):
        # 鼠标离开列表项时恢复原色
        for i in range(self.listbox.size()):
            self.listbox.itemconfig(i, bg='white')

if __name__ == "__main__":
    root = tk.Tk()
    app = GatewayScannerApp(root)
    root.mainloop()