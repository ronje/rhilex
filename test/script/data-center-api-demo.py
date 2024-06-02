import tkinter as tk
import requests


def fetch_data():
    url = "http://192.168.10.187:2580/api/v1/datacenter/queryLastData?uuid=SCHEMAZRJ5MLZN&secret=rhilex-secret"
    response = requests.get(url)
    if response.status_code == 200:
        data = response.json().get("data", {})
        update_ui(data)
    else:
        print("Failed to fetch data")


def update_ui(data):
    conductivity_label.config(text=f"Conductivity: {data.get('conductivity', '-')}")
    ph_value_label.config(text=f"pH Value: {data.get('ph_value', '-')}")
    resistivity_label.config(text=f"Resistivity: {data.get('resistivity', '-')}")
    temp_label.config(text=f"Temperature: {data.get('temp', '-')}")
    last_update_label.config(text=f"Last Update: {data.get('create_at', '-')}")

    root.after(5000, fetch_data)  # 5秒后再次获取数据


root = tk.Tk()
root.title("Data Display")

container = tk.Frame(root, bg="white", padx=20, pady=20)
container.pack(expand=True, fill="both")

conductivity_label = tk.Label(container, text="Conductivity: -")
ph_value_label = tk.Label(container, text="pH Value: -")
resistivity_label = tk.Label(container, text="Resistivity: -")
temp_label = tk.Label(container, text="Temperature: -")
last_update_label = tk.Label(container, text="Last Update: -")

conductivity_label.pack()
ph_value_label.pack()
resistivity_label.pack()
temp_label.pack()
last_update_label.pack()

fetch_data()  # 初始化获取数据

root.mainloop()
