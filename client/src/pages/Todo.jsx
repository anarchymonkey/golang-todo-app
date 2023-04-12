import React, { useEffect, useReducer, useState } from "react"

// hooks
import { useFetch } from "../hooks"

// others
import * as config from "../config.json"

// icons

import DeleteIcon from "../../public/assets/trash.svg"
import ClockIcon from "../../public/assets/clockIcon.svg"
import EditIcon from "../../public/assets/editIcon.svg"

// styles

import style from "./todo.module.css"
import { Modal } from "../views/Shared/Modal"


const handleCollectionStates = (state, action) => {
	console.log(state, action)
}


const Collections = () => {
	const [state, dispatch] = useReducer(handleCollectionStates, {
		groups: [],
		items: [],
		selectedGroup: {},
		selectedItem: {},
		contents: [],
		selectedContents: [],
	})
	const [groups, setGroups] = useState([])
	const [items, setItems] = useState([])
	const [selectedGroup, setSelectedGroup] = useState(null)
	const [selectedItem, setSelectedItem] = useState(null)
	const [contents, setContents] = useState([])
	const [selectedContents, setSelectedContents] = useState([])
	const [content, setContent] = useState("")

	const { get, loading, error, deleteRequest, post } = useFetch()

	console.log({ loading })

	useEffect(() => {
		get(config.url.getGroups).then(resp => {
			setGroups(resp)
		})
	}, [])

	const fetchItems = async (groupId) => {
		return get(`http://localhost:8080/group/${groupId}/items`)
	}

	const fetchContents = async (itemId) => {
		return get(`http://localhost:8080/item/${itemId}/contents`)
	}


	const onGroupClick = async (group) => {
		setSelectedGroup(group)
		const items = await fetchItems(group.id)
		setItems(items)
	}

	const onItemClick = async (item) => {
		setSelectedItem(item)
		const resp = await fetchContents(item.id)
		console.log({ resp })
		setContents(resp ? resp : [])
	}

	const deleteItem = async (groupId, itemId) => {
		await deleteRequest(`http://localhost:8080/group/${groupId}/item/${itemId}/delete`).then(async (res) => {
			console.log(res)
			const items = await fetchItems(groupId)
			setItems(items)
		})
	}

	const onContentChange = (event) => {
		const { value } = event.target
		setContent(value)
	}

	const onContentAddClick = async (item) => {
		console.log({ content })
		const resp = await post(`http://localhost:8080/item/${item.id}/content/add`, {
			content,
		})

		console.log({ addContentResponse: resp })
		const contentList = await fetchContents(item.id)
		setContents(contentList)
	}

	const onCheckboxChange = (content) => {
		console.log({ content, selectedContents })

		if (selectedContents.includes(content.id)) {
			const filteredSelectedContents = selectedContents.filter(
				contentId => contentId !== content.id)			
			setSelectedContents(filteredSelectedContents)
			return
		}

		setSelectedContents((previouslySelectedContents) => [
			...previouslySelectedContents, content.id])
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
			{selectedItem && (<Modal>
				<div className={style.mainContents}>
					<div>
						<input type="text" onChange={onContentChange} value={content} placeholder="Add content..." />
						<button onClick={() => onContentAddClick(selectedItem)}>Add</button>
					</div>
					<div>
						{contents.map(content => (
							<div className={style.todoContainer}  key={content.id}>
								<div className={style.contentContainer}>
									<input 
										type="checkbox" 
										checked={selectedContents.includes(content.id)} 
										onChange={() =>onCheckboxChange(content)} 
									/>
									<span>{content.content}</span>
								</div>
								<EditIcon width={20} height={20} />
							</div>
						))}
					</div>
				</div>
			</Modal>
			)}
		</div>
	)
}

export default Collections