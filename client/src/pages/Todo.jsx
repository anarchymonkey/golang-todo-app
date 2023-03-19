import { useEffect } from "react";
import { useFetch } from "../hooks/useFetch";

import * as config from '../config.json';


const Todo = () => {
    const {
        get
    } = useFetch();

    const fetchData = async () => {
        try {
            const resp = await get(config.url.getTodos);
            console.log("resp", resp);
        } catch (err) {
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
        </>
    )
}

export default Todo;