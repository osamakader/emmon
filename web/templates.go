package web

// GetHTML returns the HTML content for the web interface
func GetHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Embedded Linux Monitor</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Courier New', monospace;
            background: #1a1a1a;
            color: #00ff00;
            padding: 20px;
            font-size: 14px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .header {
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 2px solid #00ff00;
            padding-bottom: 10px;
        }
        
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .card {
            background: #2a2a2a;
            border: 1px solid #00ff00;
            border-radius: 5px;
            padding: 20px;
            box-shadow: 0 0 10px rgba(0, 255, 0, 0.3);
        }
        
        .card h3 {
            margin-bottom: 15px;
            color: #00ff00;
            border-bottom: 1px solid #00ff00;
            padding-bottom: 5px;
        }
        
        .metric {
            display: flex;
            justify-content: space-between;
            margin-bottom: 8px;
            padding: 5px 0;
        }
        
        .metric:nth-child(even) {
            background: rgba(0, 255, 0, 0.1);
        }
        
        .progress-bar {
            width: 100%;
            height: 20px;
            background: #1a1a1a;
            border: 1px solid #00ff00;
            border-radius: 3px;
            overflow: hidden;
            margin-top: 5px;
        }
        
        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #00ff00, #00cc00);
            transition: width 0.3s ease;
        }
        
        .status {
            text-align: center;
            padding: 10px;
            margin-bottom: 20px;
            border-radius: 5px;
        }
        
        .status.connected {
            background: rgba(0, 255, 0, 0.2);
            border: 1px solid #00ff00;
        }
        
        .status.disconnected {
            background: rgba(255, 0, 0, 0.2);
            border: 1px solid #ff0000;
            color: #ff0000;
        }
        
        .gpio-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
            gap: 10px;
        }
        
        .gpio-pin {
            text-align: center;
            padding: 10px;
            border: 1px solid #00ff00;
            border-radius: 3px;
            background: #1a1a1a;
        }
        
        .gpio-pin.active {
            background: #00ff00;
            color: #000;
        }
        
        @media (max-width: 768px) {
            body {
                font-size: 12px;
                padding: 10px;
            }
            
            .grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧠 Embedded Linux Monitor</h1>
            <div id="status" class="status disconnected">Disconnected</div>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3>CPU</h3>
                <div class="metric">
                    <span>Usage:</span>
                    <span id="cpu-usage">--</span>
                </div>
                <div class="progress-bar">
                    <div id="cpu-progress" class="progress-fill" style="width: 0%"></div>
                </div>
                <div class="metric">
                    <span>Load Average (1m):</span>
                    <span id="cpu-load-1">--</span>
                </div>
                <div class="metric">
                    <span>Load Average (5m):</span>
                    <span id="cpu-load-5">--</span>
                </div>
                <div class="metric">
                    <span>Load Average (15m):</span>
                    <span id="cpu-load-15">--</span>
                </div>
                <div class="metric">
                    <span>Frequency:</span>
                    <span id="cpu-freq">--</span>
                </div>
            </div>
            
            <div class="card">
                <h3>Memory</h3>
                <div class="metric">
                    <span>Usage:</span>
                    <span id="mem-usage">--</span>
                </div>
                <div class="progress-bar">
                    <div id="mem-progress" class="progress-fill" style="width: 0%"></div>
                </div>
                <div class="metric">
                    <span>Total:</span>
                    <span id="mem-total">--</span>
                </div>
                <div class="metric">
                    <span>Used:</span>
                    <span id="mem-used">--</span>
                </div>
                <div class="metric">
                    <span>Free:</span>
                    <span id="mem-free">--</span>
                </div>
                <div class="metric">
                    <span>Available:</span>
                    <span id="mem-available">--</span>
                </div>
            </div>
            
            <div class="card">
                <h3>Disk</h3>
                <div class="metric">
                    <span>Usage:</span>
                    <span id="disk-usage">--</span>
                </div>
                <div class="progress-bar">
                    <div id="disk-progress" class="progress-fill" style="width: 0%"></div>
                </div>
                <div class="metric">
                    <span>Total:</span>
                    <span id="disk-total">--</span>
                </div>
                <div class="metric">
                    <span>Used:</span>
                    <span id="disk-used">--</span>
                </div>
                <div class="metric">
                    <span>Free:</span>
                    <span id="disk-free">--</span>
                </div>
                <div class="metric">
                    <span>I/O Read:</span>
                    <span id="disk-io-read">--</span>
                </div>
                <div class="metric">
                    <span>I/O Write:</span>
                    <span id="disk-io-write">--</span>
                </div>
            </div>
            
            <div class="card">
                <h3>Temperature</h3>
                <div class="metric">
                    <span>CPU:</span>
                    <span id="temp-cpu">--</span>
                </div>
                <div class="metric">
                    <span>GPU:</span>
                    <span id="temp-gpu">--</span>
                </div>
                <div class="metric">
                    <span>Board:</span>
                    <span id="temp-board">--</span>
                </div>
                <div class="metric">
                    <span>Ambient:</span>
                    <span id="temp-ambient">--</span>
                </div>
            </div>
        </div>
        
        <div class="card">
            <h3>GPIO Status</h3>
            <div id="gpio-container" class="gpio-grid">
                <div class="gpio-pin">No GPIO data</div>
            </div>
        </div>
    </div>

    <script>
        let ws = null;
        let reconnectTimer = null;
        
        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                document.getElementById('status').textContent = 'Connected';
                document.getElementById('status').className = 'status connected';
                if (reconnectTimer) {
                    clearTimeout(reconnectTimer);
                    reconnectTimer = null;
                }
            };
            
            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                updateDisplay(data);
            };
            
            ws.onclose = function() {
                document.getElementById('status').textContent = 'Disconnected - Reconnecting...';
                document.getElementById('status').className = 'status disconnected';
                
                if (!reconnectTimer) {
                    reconnectTimer = setTimeout(connect, 3000);
                }
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        function updateDisplay(data) {
            // Update CPU
            document.getElementById('cpu-usage').textContent = data.cpu.usage_percent.toFixed(1) + '%';
            document.getElementById('cpu-progress').style.width = data.cpu.usage_percent + '%';
            
            if (data.cpu.load_average && data.cpu.load_average.length >= 3) {
                document.getElementById('cpu-load-1').textContent = data.cpu.load_average[0].toFixed(2);
                document.getElementById('cpu-load-5').textContent = data.cpu.load_average[1].toFixed(2);
                document.getElementById('cpu-load-15').textContent = data.cpu.load_average[2].toFixed(2);
            }
            
            document.getElementById('cpu-freq').textContent = (data.cpu.frequency / 1000).toFixed(1) + ' GHz';
            
            // Update Memory
            document.getElementById('mem-usage').textContent = data.memory.usage_percent.toFixed(1) + '%';
            document.getElementById('mem-progress').style.width = data.memory.usage_percent + '%';
            document.getElementById('mem-total').textContent = formatBytes(data.memory.total);
            document.getElementById('mem-used').textContent = formatBytes(data.memory.used);
            document.getElementById('mem-free').textContent = formatBytes(data.memory.free);
            document.getElementById('mem-available').textContent = formatBytes(data.memory.available);
            
            // Update Disk
            document.getElementById('disk-usage').textContent = data.disk.usage_percent.toFixed(1) + '%';
            document.getElementById('disk-progress').style.width = data.disk.usage_percent + '%';
            document.getElementById('disk-total').textContent = formatBytes(data.disk.total);
            document.getElementById('disk-used').textContent = formatBytes(data.disk.used);
            document.getElementById('disk-free').textContent = formatBytes(data.disk.free);
            document.getElementById('disk-io-read').textContent = formatBytes(data.disk.io_read);
            document.getElementById('disk-io-write').textContent = formatBytes(data.disk.io_write);
            
            // Update Temperature
            document.getElementById('temp-cpu').textContent = data.temperature.cpu.toFixed(1) + '°C';
            document.getElementById('temp-gpu').textContent = data.temperature.gpu.toFixed(1) + '°C';
            document.getElementById('temp-board').textContent = data.temperature.board.toFixed(1) + '°C';
            document.getElementById('temp-ambient').textContent = data.temperature.ambient.toFixed(1) + '°C';
            
            // Update GPIO
            updateGPIO(data.gpio.pins);
        }
        
        function updateGPIO(pins) {
            const container = document.getElementById('gpio-container');
            container.innerHTML = '';
            
            if (!pins || Object.keys(pins).length === 0) {
                container.innerHTML = '<div class="gpio-pin">No GPIO data</div>';
                return;
            }
            
            for (const [pinName, pinData] of Object.entries(pins)) {
                const pinElement = document.createElement('div');
                pinElement.className = 'gpio-pin' + (pinData.value === 1 ? ' active' : '');
                pinElement.innerHTML = '<div><strong>' + pinName + '</strong></div><div>' + pinData.value + '</div><div>' + pinData.mode + '</div>';
                container.appendChild(pinElement);
            }
        }
        
        function formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }
        
        // Connect on page load
        connect();
    </script>
</body>
</html>`
}
