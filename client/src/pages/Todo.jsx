import React, { useEffect, useState } from "react"

// hooks
import { useFetch } from "../hooks"

// others
import * as config from "../config.json"

// styles

import style from "./todo.module.css"


const Collections = () => {

	const [groups, setGroups] = useState([])
	const [items, setItems] = useState([])

	const { get, loading, error } = useFetch()

	console.log(loading, error)

	useEffect(() => {
		get(config.url.getGroups).then(resp => {
			setGroups(resp)
		})
	}, [])

	const onGroupClick = async (group) => {
		const items = await get(`http://localhost:8080/group/${group.id}/items`);
		setItems(items)
	}

	const onItemClick = async (item) => {
		const resp = await get(`http://localhost:8080/item/${item.id}/contents`)
		console.log({ resp })
	}

	return (
		<div className={style.main}>
			<div className={style.groups}>
				<button style={{ margin: "20px 0px" }} className={style.addGroup} type="button">
					Add group
				</button>
				<div className={style.groupListItems}>
					{!error && groups.map(group => (
						<section key={group.id} onClick={() => onGroupClick(group)}>{group.title}</section>
					))}
				</div>
			</div>
			<div className={style.contents}>
				<div className={style.contentHeader}>
					<button className={style.addGroup} type="button">
						Add item
					</button>
					<div>
						<input type="search" placeholder="Search items" />
					</div>
				</div>
				<div className={style.contentItemGrid}>
					{!error && items?.map(item => (
						<section key={item.id} onClick={() => onItemClick(item)}>
							<span className={style.topLeft}>
								{new Date(item.created_at).toDateString()}
							</span>
							<span>
								{item.content}
							</span>
						</section>
					))}
				</div>
			</div>
		</div>
	)
}

export default Collections