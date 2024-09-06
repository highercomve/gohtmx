![image of the example website](https://github.com/highercomve/gohtmx/blob/main/full.png?raw=true)

# Golang and HTMX combination for embedded linux

Golang and HTMX combination for embedded linux

For the past six years, I’ve been working at Pantacor, a company that develops an open-source framework for containerized embedded Linux. Our focus is on creating embedded Linux solutions using container-based building blocks. As a web developer, my experience has primarily been in building architectures for web applications and APIs. Embedded systems wasn’t part of my toolkit since university. However, the need to create small applications with web components for embedded Linux led me to explore Go and HTMX as an alternative to our existing Rust + React applications or shell API and React combinations.

## The Challenge of Embedded Linux Development

Embedded Linux systems present unique challenges, particularly when developing user interfaces and web-based applications. These systems often have constrained resources — limited memory and processing power — making it crucial to choose efficient and lightweight technologies.

## Exploring Golang and HTMX

Golang, or Go, is a statically typed, compiled language known for its simplicity and efficiency. HTMX is a lightweight JavaScript library that allows access to modern browser features directly from HTML, reducing the need for complex JavaScript frameworks and eliminating the need for a whole transpilation and compilation step to build your application.

In our embedded Linux context, applications need to be self-contained, including all assets that need to be served offline. This means creating an application that can be served from the device’s hotspot or if you are connected directly via ethernet.

The combination of Golang and HTMX offers several advantages for our embedded Linux applications:

* Resource Efficiency: Both Go and HTMX have small footprints, making them suitable for resource-constrained embedded systems. In most cases, we need to create a container with the whole application without exceeding 5MB for the entire package.

* Approachable Learning Curve: Given my prior experience with Go in our backends and CLI applications, the learning curve primarily revolves around the specific problem I need to solve, rather than the language itself.

* Performance: Go’s compiled nature and HTMX’s minimal JavaScript approach result in efficient execution and reduced overhead, crucial for our limited-resource environments.

* Offline Capability: The ability to package all necessary assets within the application aligns perfectly with our need for offline functionality in embedded systems.

* Simplified Development Process: By eliminating complex build processes and reducing the amount of JavaScript needed, we can streamline our development and deployment workflows.

## Real-World Applications

To illustrate the practical benefits of using Go and HTMX in embedded Linux environments, let’s explore two examples from our work at Pantacor:

## 1. WiFi Connect App

One of the applications we’ve developed is a WiFi connection tool for Linux devices. This app creates a captive portal, allowing users to set up WiFi on their embedded Linux device easily. Here’s how it works:

* The device creates an access point.

* Users connect to this access point using their phone or laptop.

* Through a web interface, users can provide WiFi credentials for their home network.

* The device then uses these credentials to connect to the user’s home network.

## 2. Device Management App

Another example is our device management application for Pantavisor. This web-based utility allows users to manage their Linux devices through a browser interface. Features include:

* Installing new containers by drag-and-dropping tarball files

* Uninstalling containers with a few clicks

* Reading container logs

* Setting up container configurations

## Implementing Go and HTMX: A Practical Example

Let’s dive into a practical example of how to use Go and HTMX to create a simple WiFi network management application. This example will demonstrate how to list available networks and connect to a selected network.

The source code of the whole example can be found here:

[https://github.com/highercomve/gohtmx](https://github.com/highercomve/gohtmx)

## Setting Up the Server

First, we’ll create a small server with Go that has three endpoints:

1. GET /: Returns the index.html of the application

1. GET /networks: Returns the list of networks as HTML

1. POST /networks: Saves the selected network and connects to it

Here’s how we implement the index endpoint:

    func getIndex(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
     data := servermodels.Response{
      Code:  http.StatusOK,
      Data:  nil,
      Error: nil,
     }
     render(w, r, opts, "index.html", &data)
    }

## The Index HTML

The index.html file uses HTMX to create an interactive interface:

    {{template "layout_start"}}
    {{define "title"}}Go HTMX Example{{end}}
    <section class="flex justify-content-center pb-2">
     <header class="col-6 text-center">
      <h1>Manage wifi networks with Golang + HTMX + nmcli</h1>
     </header>
    </section>
    <section class="flex justify-content-center mb-2">
     <header class="col-12 flex align-items-center justify-content-center mb-2">
      <button
       class="btn btn-primary"
       hx-get="/networks"
       hx-target="#connections-list"
       hx-swap="morph:outerHTML"
      >
       Scan Networks
       <div class="htmx-indicator spinner-grow" role="status"></div>
      </button>
     </header>
     {{if notNil .Data}} {{ block "network_list" . }}{{end}} {{else}}
     <section id="connections-list">
      <div
       hx-get="/networks"
       hx-target="#connections-list"
       hx-swap="morph:outerHTML"
       hx-trigger="load"
       class="flex flex-direction-column justify-content-center align-items-center"
      >
       <div class="htmx-indicator spinner-grow" role="status"></div>
       <div class="pt-1 color-grey-300">Loading networks...</div>
      </div>
     </section>
     {{end}}
    </section>
    {{template "layout_end"}}

## Understanding HTMX Attributes

Let’s take a closer look at the HTMX attributes we’re using in our WiFi management application. These attributes are what give our application its dynamic, app-like feel without the need for complex JavaScript:

1. hx-get="/networks": This attribute tells HTMX to make a GET request to the "/networks" endpoint when the element is triggered. In our "Scan Networks" button, this means clicking the button will fetch the list of available networks.

1. hx-target="#connections-list": This specifies where the response from the server should be inserted. In our case, the list of networks will be placed into the element with the ID "connections-list".

1. hx-swap="morph:outerHTML": This determines how the new content is swapped in. The "morph" strategy intelligently updates the DOM, minimizing changes for a smooth transition. "outerHTML" means it replaces the entire target element, including the element itself.

1. hx-trigger="load": This attribute defines when the HTMX request should be triggered. When set to "load", it means the request will be made as soon as the element is loaded into the DOM. We use this in our initial network list div to automatically load the networks when the page loads.

1. <div class="htmx-indicator spinner-grow" role="status"></div>: While not an hx-* attribute, this div is part of HTMX's built-in system for showing loading indicators. The "htmx-indicator" class tells HTMX to show this element while a request is in progress and hide it otherwise.

These HTMX attributes allow us to create a responsive, dynamic user interface with minimal JavaScript. The server sends HTML in response to these requests, which HTMX then seamlessly integrates into the page. This approach is particularly beneficial in our embedded Linux context, where we want to keep client-side processing to a minimum while still providing a smooth user experience.

By leveraging these HTMX features, we can create an interactive WiFi management interface that feels responsive and modern, all while keeping our application lightweight and efficient — crucial factors when working with embedded Linux systems.

## Listing Available Networks

The GET /networks endpoint retrieves and displays the list of available WiFi networks:

    func getNetworksList(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
      networkmanager := nm.Init()
      conns, err := networkmanager.List()
      if err != nil {
        opts.Logger.Println(err)
        data := servermodels.Response{
          Code:  http.StatusInternalServerError,
          Data:  nil,
          Error: err,
        }
        render(w, r, opts, "error.html", &data)
        return
      }
      data := servermodels.Response{
        Code:  http.StatusOK,
        Data:  conns,
        Error: nil,
      }
      render(w, r, opts, "network_list", &data)
    }

The networkmanager.List() function uses the nmcli command to retrieve the list of available WiFi networks:

    func (nm *NetworkManager) List() (conns []nmmodules.WifiConn, err error) {
      cmd := exec.Command(
        "nmcli",
        "-f",
        "ssid,mode,freq,signal,active,security",
        "device",
        "wifi",
        "list",
     )
     stdout, err := cmd.Output()
     if err != nil {
       fmt.Println("Error running nmcli command:", err)
       return
     }

     // ... (code to parse the output and create WifiConn structs)
     return conns, err
    }

## **The list of network HTML**

    {{ define "network_list" }} {{ $numOfConns := len .Data }}
    <section
     id="connections-list"
     class="connections flex justify-content-center col-12"
    >
     <form
      class="col-6"
      hx-post="/networks"
      hx-target="body"
      hx-confirm="After applying this changes the device is going to be connected to the new network and the AP will disapear."
     >
      {{if notNil .Error}}
      <div class="alert alert-danger">{{.Error}}</div>
      {{end}}
      <div
       class="input-group mb-2 flex align-items-center justify-content-center"
      >
       <details class="custom-select">
        <summary class="radios">
         {{ if (eq $numOfConns 0) }}
         <input type="radio" title="No networks available" checked />
         {{ else }} {{range $conn := .Data}}
         <input
          type="radio"
          name="ssid"
          id="{{$conn.ID}}"
          title="{{$conn.SSID}} {{if $conn.Active}}*{{end}}"
          value="{{$conn.SSID}}"
          {{if
          $conn.Active}}checked{{end}}
         />
         {{ end }} {{ end }}
        </summary>
        <ul class="list">
         {{ if eq $numOfConns 0 }}
         <li>
          <label>No networks</label>
         </li>
         {{ else }} {{range $conn := .Data}}
         <li>
          <label for="{{$conn.ID}}">
           {{block "signal" $conn.Strength}}{{end}}
           {{$conn.SSID}}
           <span class="ml-05"> {{$conn.Frequency}} </span>
           {{ if $conn.Active }} * {{ end }}
          </label>
         </li>
         {{end}} {{ end }}
        </ul>
       </details>
      </div>
      <div class="input-group mb-1 flex">
       <label class="input-group-text" for="password">Password</label>
       <input
        name="password"
        id="password"
        type="password"
        aria-label="Password"
        class="form-control"
       />
      </div>
      <div class="input-group mb-2 flex">
       <input
        type="checkbox"
        id="show-password"
        class=""
        onclick="showPassword('password')"
       />
       <label class="input-group-text" for="show-password"
        >Show Password</label
       >
      </div>
      <script>
       function showPassword(elementID) {
        var x = document.getElementById(elementID);
        if (x.type === "password") {
         x.type = "text";
        } else {
         x.type = "password";
        }
       }
      </script>
      <div class="input-group mb-2 flex justify-content-center">
       <button class="btn btn-primary" type="submit">
        Connect to this network
        <div class="htmx-indicator spinner-grow" role="status"></div>
       </button>
      </div>
     </form>
    </section>
    {{ end }}

This HTML will generate an interface that looks like this:

![Image of how the network_list.html looks when is rendered by the browser](https://cdn-images-1.medium.com/max/2000/1*zvRMb7k7VOh4DxHevUMqjw.png)*Image of how the network_list.html looks when is rendered by the browser*

## Connecting to a Network

In this section, we’ll walk through how to connect to a WiFi network using nmcli, a command-line tool for controlling NetworkManager. The process will be integrated into the Golang and HTMX applications for a seamless connection experience.

### The POST /networks Endpoint

When a user selects a network from the list and submits their password, the POST /networks endpoint is responsible for connecting the device to the chosen WiFi network. Here's a breakdown of the flow:

1. **User Submits a Network and Password:** The user interface allows the user to choose a network and enter the password.

1. **Form Submission:** The form data (network SSID and password) is submitted via POST to the /networks endpoint.

1. **NetworkManager CLI (nmcli) Integration:** The Go application invokes nmcli to handle the connection to the selected WiFi network.

How to handle the network creation

    func handleNetworkConnection(w http.ResponseWriter, r *http.Request, opts *servermodels.ServerOptions) {
       ssid := r.FormValue("ssid")
       password := r.FormValue("password")
       networkmanager := nm.Init()

       // Call the function to connect using nmcli
       err := networkmanager.Save(ssid, password)
       if err != nil {
        opts.Logger.Println(err)
          data := servermodels.Response{
            Code:  http.StatusOK,
            Data:  nil,
            Error: err,
          }
          render(w, r, opts, "error.html", &data)
          return
       }

       data := servermodels.Response{
          Code:  http.StatusOK,
          Data:  nil,
          Error: nil,
       }
       opts.Logger.Printf("Connected successfully to %s\n", ssid)
       render(w, r, opts, "", &data)
    }

the save connection method will be

    func (nm *NetworkManager) Save(ssid, password string) error {
     // Use nmcli to connect to the network
     cmd := exec.Command("nmcli", "device", "wifi", "connect", ssid, "password", password)

     // Capture both the output and errors separately
     var stderr bytes.Buffer
     var stdout bytes.Buffer
     cmd.Stdout = &stdout
     cmd.Stderr = &stderr

     // Execute the command
     err := cmd.Run()
     if err != nil {
      // If there's an error, return it as a Go error, including nmcli's stderr
      return fmt.Errorf("Failed to connect to network %s: %s, nmcli error: %s", ssid, err.Error(), stderr.String())
     }

     // Optionally log the success message from stdout
     fmt.Println("nmcli output:", stdout.String())

     return nil
    }

## Conclusion

The use of Go and HTMX for embedded Linux applications has significantly streamlined our development process at Pantacor. It allows us to create efficient, self-contained applications that meet the unique challenges of embedded systems. While this approach may not be suitable for every scenario, it has proven to be a powerful tool in our toolkit, especially for devices with limited resources that require user-friendly interfaces.

As embedded systems continue to evolve and become more prevalent, the need for efficient, lightweight web technologies will only grow. The combination of Go and HTMX represents a promising direction for developers looking to build modern web interfaces for embedded Linux devices.

By leveraging Go’s efficiency and HTMX’s simplicity, we can create responsive, feature-rich applications that respect the constraints of embedded systems while providing a smooth user experience. This approach not only simplifies development but also ensures that our applications remain maintainable and adaptable to future needs.
