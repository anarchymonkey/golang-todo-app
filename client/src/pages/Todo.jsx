import { createContext, useEffect, useMemo, useState } from "react";
import { useFetch } from "../hooks/useFetch";
import { Todos } from "../views/Todo";

import * as config from '../config.json';

export const TodoContext = createContext();
const Todo = () => {
    const {
        get
    } = useFetch();

    const [todos, setTodos] = useState([]);

    const fetchData = async () => {
        try {
            const resp = await get(config.url.getTodos);
            setTodos(() => resp);
        } catch (err) {
            setTodos([]);
            console.log("The error is", err)
        }
    }


    useEffect(() => {
        fetchData()
    }, []);

    return (
        <div style={{ display: 'flex', justifyContent: "center", alignItems: "center" }}>
            <TodoContext.Provider value={{todos, setTodos, fetchData }}>
                <Todos />
            </TodoContext.Provider>
        </div>
    )
}

export default Todo;