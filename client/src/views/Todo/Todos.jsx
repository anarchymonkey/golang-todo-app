import { useMemo, useRef, useState } from "react";
import { Card } from "../Shared";
import TodoItem from "./TodoItem";

const Todos = ({
    data,
}) => {
    const checkboxRef = useRef(null);
    const [selectedItems, setSelectedItems] = useState(new Set());

    return (
        <Card>
            {data.map((todo) => (
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