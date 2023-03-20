import { useEffect, useMemo, useState } from "react";
import { useFetch } from "../hooks/useFetch";
import { Todos } from "../views/Todo";

import * as config from '../config.json';


const Todo = () => {
    const {
        get
    } = useFetch();
    
    const [todos, setTodos] = useState([]);

    const fetchData = async () => {
        try {
            const resp = await get(config.url.getTodos);
            setTodos((prevResp) => resp);
        } catch (err) {
            setTodos([]);
            console.log("The error is", err)
        }
    }


    useEffect(() => {
        fetchData()
    }, []);

    return (
        <>
            <div>Header</div>
            <div>Subtitle</div>
            <div style={{ display: 'flex', justifyContent: "center", alignItems: "center"}}>
                <Todos data={todos} />
            </div>
        </>
    )
}

export default Todo;