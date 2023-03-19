
const TodoItem = ({ todoItem }) => {
    return (
        <div style={{
            display: 'flex',
            justifyContent: "flex-start",
            alignItems: "center",
            gap: "10px"
        }}>
            <input type="checkbox" checked></input>
            <span>{todoItem.name}</span>
        </div>
    )
}

export default TodoItem;