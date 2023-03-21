import { useContext, useRef, useState } from "react";
import { Card } from "../Shared";
import TodoItem from "./TodoItem";

import "./styles.css";
import { TodoContext } from "../../pages/Todo";

const Todos = () => {
    const { todos, addTodo } = useContext(TodoContext);
    
    const checkboxRef = useRef(null);
    const inputAddRef = useRef(null);

    const [selectedItems, setSelectedItems] = useState(new Set());

    const onKeyPress = (event) => {
        event.preventDefault();

        if (event.which === 13 && event.code === "Enter") {
            console.log(inputAddRef.current.value);
            addTodo({
                name: inputAddRef.current.value,
            })

            inputAddRef.current.value = '';
        }

    }

    return (
        <Card>
            <input type="text" ref={inputAddRef} onKeyUp={onKeyPress}/>
            {todos.map((todo) => (
                <div key={todo.id}>
                    <TodoItem
                        todoItem={todo}
                        ref={checkboxRef}
                        checked={selectedItems.has(todo.id)}
                        setSelectedItems={setSelectedItems}
                    />
                </div>
            ))}
        </Card>
    )
}

export default Todos;