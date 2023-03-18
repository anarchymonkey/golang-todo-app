import { Profiler } from "react";
import { Layout, Todo } from "./pages";


// Profiler to check if my app is not shitting around anymore, a simple todo app does not need it but why not

const App = () => {
    return (
        <Profiler>
            <Layout>
                <Todo />
            </Layout>
        </Profiler>
    )
}

export default App;