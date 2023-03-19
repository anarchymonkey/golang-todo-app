import { Card } from "../Shared";
import TodoItem from "./TodoItem";

const Todos = ({
    data,
}) => {
    return (
        <Card>
            {data.map((todo) => (
                <div key={todo.id}>
                    <TodoItem todoItem={todo} />
                </div>
            ))}
        </Card>
    )
}

export default Todos;