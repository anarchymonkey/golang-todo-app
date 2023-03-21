import { createContext, useEffect, useMemo, useState } from "react";
import { useFetch } from "../hooks/useFetch";
import { Todos } from "../views/Todo";

import * as config from '../config.json';

export const TodoContext = createContext();
const Todo = () => {
    const {
        get,
        put,
        post,
    } = useFetch();

    const [todos, setTodos] = useState([]);

    const getTodos = async () => {
        try {
            const resp = await get(config.url.getTodos);
            setTodos(() => resp);
        } catch (err) {
            setTodos([]);
            console.log("The error is", err)
        }
    }

    const addTodo = (todoItem) => {
        post(config.url.addTodo, {
            name: todoItem.name,
        }).then(() => getTodos()); 
    }

    const updateTodo = (id, todoItem) => {
        put(config.url.updateTodo, {
            id,
            name: todoItem.name,
            isComplete: todoItem.isComplete,
            isStriked: todoItem.isStriked,
        }).then(() => {
            getTodos();
        })
    }


    useEffect(() => {
        getTodos();
    }, []);

    return (
        <div style={{ display: 'flex', justifyContent: "center", alignItems: "center" }}>
            <TodoContext.Provider value={{todos, setTodos, getTodos, updateTodo, addTodo }}>
                <Todos />
            </TodoContext.Provider>
        </div>
    )
}

export default Todo;