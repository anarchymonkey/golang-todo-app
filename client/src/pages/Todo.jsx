import React, { useEffect, useState } from "react"

// hooks
import { useFetch } from "../hooks"

// others
import * as config from "../config.json"

// icons

import DeleteIcon from "../../public/assets/trash.svg"
import ClockIcon from "../../public/assets/clockIcon.svg"

// styles

import style from "./todo.module.css"
import { Modal } from "../views/Shared/Modal"


const Collections = () => {

	const [groups, setGroups] = useState([])
	const [items, setItems] = useState([])
	const [selectedGroup, setSelectedGroup] = useState(null)

	const { get, loading, error, deleteRequest } = useFetch()

	console.log(loading, error)

	useEffect(() => {
		get(config.url.getGroups).then(resp => {
			setGroups(resp)
		})
	}, [])

	const fetchItems = async (groupId) => {
		return get(`http://localhost:8080/group/${groupId}/items`)
	}


	const onGroupClick = async (group) => {
		setSelectedGroup(group)
		const items = await fetchItems(group.id)
		setItems(items)
	}

	const onItemClick = async (item) => {
		const resp = await get(`http://localhost:8080/item/${item.id}/contents`)
		console.log({ resp })
	}

	const deleteItem = async (groupId, itemId) => {
		await deleteRequest(`http://localhost:8080/group/${groupId}/item/${itemId}/delete`).then(async (res) => {
			console.log(res)
			const items = await fetchItems(groupId)
			setItems(items)
		})
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
							<DeleteIcon className={style.deleteBtn} onClick={() => deleteItem(selectedGroup.id, item.id)} />
							{item?.remind_at && (
								<div>
									<ClockIcon className={style.clockIcon} />
									{new Date(item?.remind_at)}
								</div>
							)}
						</section>
					))}
				</div>
			</div>
			<Modal>
				<h1>Hey my name is aniket</h1>
				<h2>This would be a todo list</h2>
			</Modal>
		</div>
	)
}

export default Collections