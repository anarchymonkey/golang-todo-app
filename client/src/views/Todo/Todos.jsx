import { useContext, useRef, useState } from "react";
import { Card } from "../Shared";
import TodoItem from "./TodoItem";

import AddCircleIcon from "../../../public/assets/AddCircleIcon.svg";
import { default as consts } from './const';

import "./styles.css";
import { TodoContext } from "../../pages/Todo";

const Todos = () => {
    const { todos } = useContext(TodoContext);
    const checkboxRef = useRef(null);
    const [selectedItems, setSelectedItems] = useState(new Set());

    const onAddTodoClick = () => {
        console.log("Add todo was clicked")
    }

    return (
        <Card>
            <AddCircleIcon
                width={50}
                height={50}
                id="add-icon"
                onClick={onAddTodoClick}
            />
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