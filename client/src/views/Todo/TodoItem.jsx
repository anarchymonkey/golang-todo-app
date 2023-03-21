import { forwardRef, useContext, memo, useMemo, useState } from "react";
import EditIcon from "../../../public/assets/editIcon.svg"
import { default as consts } from './const.js';
import * as config from '../../config.json';
import { TodoContext } from "../../pages/Todo";


import './styles.css';
import { useFetch } from "../../hooks/useFetch";

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
    const { updateTodo } = useContext(TodoContext);
    const { put } = useFetch();
    const [mode, setMode] = useState(consts.modes.DEFAULT);
    const [currentContent, setCurrentContent] = useState(todoItem.name);

    const onTextContentChange = (event) => {
        if (event.target) {
            console.log(event);
            const { value } = event.target;
            setCurrentContent(value);
        }
    }


    const getTodoListMode = useMemo(() => {
        if (mode === consts.modes.UPDATE) {
            return <input type="text" id="update-todo" value={currentContent} onChange={onTextContentChange} />
        }
        return currentContent;
    }, [mode, currentContent]);
    

    const onTodoUpdateCompletion = () => {
        updateTodo(todoItem.id, {
            name: currentContent,
            isComplete: false,
            isStriked: false,
        });

        setMode(consts.modes.DEFAULT);

        
    };


    return (
        <div className="todo-item">
            <div className="todo-item-children left">
                <input
                    type="checkbox"
                    ref={ref}
                    checked={checked}
                    id={todoItem.id}
                    onChange={(event) => onChange(event, setSelectedItems)} />
                <span>{getTodoListMode}</span>
            </div>
            <div className="todo-item-children right">
                <span>
                    {mode === consts.modes.UPDATE ? (
                        <div onClick={onTodoUpdateCompletion}>Done</div>
                    ) : (
                        <EditIcon
                            width={20}
                            height={20}
                            fill="white"
                            onClick={() => setMode(consts.modes.UPDATE)}
                        />
                    )}
                </span>
            </div>
        </div>
    )
});

export default memo(TodoItem);