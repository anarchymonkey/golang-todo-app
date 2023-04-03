import { Profiler } from "react";
import { Header, Layout, Todo } from "./pages";
import { BrowserRouter } from 'react-router-dom';


// Profiler to check if my app is not shitting around anymore, a simple todo app does not need it but why not

const App = () => {
    return (
        <Profiler>
            <BrowserRouter>
                <Layout>
                    <Header />
                    <Todo />
                </Layout>
            </BrowserRouter>
        </Profiler>
    )
}

export default App;