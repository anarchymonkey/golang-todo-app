import { forwardRef, memo, useMemo } from "react";


// what debounce does
// takes a function and a delay
// calls the function after a delay, if the function is triggered again it starts the timer again
const debounce = (callbackFunction, delay) => {

}
const onChange = (event, setSelectedItems) => {
    const { id } = event.target;
    setSelectedItems((prevSet) => {
        const items = new Set(prevSet)

        if (items.has(id)) {
            items.delete(id)
            return items;
        }
        items.add(id)
        return items;
    })
}

const TodoItem = forwardRef(({ todoItem, checked, setSelectedItems }, ref) => {    

    console.log({ checked });
    
    return (
        <div style={{
            display: 'flex',
            justifyContent: "flex-start",
            alignItems: "center",
            gap: "10px"
        }}>
            <input
                type="checkbox"
                ref={ref}
                checked={checked}
                id={todoItem.id}
                onChange={(event) => onChange(event, setSelectedItems)} />
            <span>{todoItem.name}</span>
        </div>
    )
});

export default memo(TodoItem);