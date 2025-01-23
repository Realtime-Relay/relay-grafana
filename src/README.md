# Relay
<b>"Monitor the Unmonitorable: Real-Time Insights for Your Mission-Critical Systems"</b><br>
The Grafana Real-Time Data Source Plugin is designed for mission-critical observability, enabling you to monitor and act on live data from cyber-physical systems like rockets, rovers, EV charging stations, and more. With seamless integration with Node.js and Python libraries, flexible messaging APIs, and real-time alerting, this tool empowers you to customize and scale your monitoring to meet the demands of your unique applications.

## Features
1. <b>Custom Data Structures:</b> Visualize user-defined data structures, allowing complete flexibility to adapt to your system's unique requirements.
2. <b>Real-Time Alerting:</b> Set up instant alerts based on live data streams to ensure you're always informed of critical events.
3. <b>Seamless Integration:</b> Publish data effortlessly using Node.js and Python libraries, making it easy to integrate with your existing workflows.
4. <b>Low-Latency Streaming:</b> Achieve near-zero delay in data visualization, perfect for monitoring fast-moving systems like rockets, rovers, and EVs.
5. <b>Real-Time Dashboards:</b> Create interactive Grafana dashboards for live data visualization tailored to your operational needs.
6. <b>Configurable Alerting Rules:</b> Define custom alert thresholds and conditions to receive actionable insights from your data streams.
7. <b>Cyber-Physical System Ready:</b> Specifically designed for real-time observability of complex systems like IoT devices, EV infrastructure, and space missions.

## Version Compatibility
Minimum supported Grafana version => 10.4.0

## Setup
1. Obtain API key and secret key
2. Enter the API Key & Secret into their respective fields.
3. In the 'Path' field, enter `api.relay-x.io`
4. Click on the `Save & Test` button and expect a `Connection successfully established` message in green. If you get this, the data source is now setup.
![Datasource Setup](https://github.com/Realtime-Relay/relay-grafana/releases/download/v0.0.1/ds_setup.png "Datasource Setup")

## Example
To demonstrate the how the data source works, a [script](https://github.com/Realtime-Relay/relayx-js/blob/main/examples/example_send_data_on_connect.js) will publish to a topic called "power-telemetry". The data source listens to "power-telemetry" and displays it on a time series graph.<br>

The script generates random values between 0 and 100, sends it to the Relay Network and relays it to the datasource on grafana.

### Dashboard Setup
1. Create a Time Series Panel
   * Add a new time series graph to your Grafana dashboard.
   * Select the Relay Data Source as the data source for the panel.
2. Set the Query Topic
   * In the query editor, specify the topic as "power-telemetry". This will fetch data streamed to that topic.
3. Define the Time Range
   * Configure the time range to now-6s to now for a scrolling, real-time graph. While this is an example range, keep it within 10s.
4. Transform the Data
   * Apply necessary transformations to structure the incoming data from the topic according to the visualization requirements of your panel.
5. Start Streaming
   * Hit the Refresh button. As long as data is being sent to the "power-telemetry" topic, it will populate the graph in real time.
<br>

![Graph Setup](https://github.com/Realtime-Relay/relay-grafana/releases/download/v0.0.1/graph_setup.png "Graph Setup")

### Demo
![Demo](https://github.com/Realtime-Relay/relay-grafana/releases/download/v0.0.1/demo_gif.gif "Demo")