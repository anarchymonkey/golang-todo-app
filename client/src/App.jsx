import React, { Profiler } from "react"
import { BrowserRouter } from "react-router-dom"

// components
import { Layout } from "./pages"


// Profiler to check if my app is not shitting around anymore, a simple todo app does not need it but why not

const App = () => {
	return (
		<Profiler>
			<BrowserRouter>
				<Layout />
			</BrowserRouter>
		</Profiler>
	)
}

export default App