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

const debounce = (fn, delay) => {
	let timer
	
	return function () {
		clearTimeout(timer)
		timer = setTimeout(() => {
			fn(...arguments)
		}, delay)
	}
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
	const [editableContents, setEditableContents] = useState(new Map())

	const { get, loading, error, deleteRequest, post, put } = useFetch()

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
		if (resp?.length === 0) {
			setContents([])
			setEditableContents(new Map())
			return
		}
		setEditableContents(() => new Map(resp?.map(res => [res.id, {
			id: res.id,
			originalValue: res.content,
			updatedValue: res.content,
			error: {
				isError: false,
				message: "",
			}
		}])))
		setContents(resp)
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

	const onContentEdit = async (ev, content) => {
		console.log({ upVal: editableContents, ev: ev, content })
		const { value } = ev.target
		if (editableContents.get(content.id).originalValue === value) {
			const updatedContents = new Map(editableContents)
			updatedContents.set(content.id, {
				...editableContents.get(content.id),
				updatedValue: value,
			})
			setEditableContents(updatedContents)
			return
		}

		const updatedContents = new Map(editableContents)
		updatedContents.set(content.id, {
			...editableContents.get(content.id),
			updatedValue: value,
		})
		setEditableContents(updatedContents)

		const debouncedFn = debounce(async () => {
			const resp = await put(`http://localhost:8080/content/${content.id}/update`, {
				content: value,
			})
	
			console.log({ resp })
		}, 1500)
	
		debouncedFn()
	}

	const onContentAddClick = async (item) => {
		console.log({ content })
		const resp = await post(`http://localhost:8080/item/${item.id}/content/add`, {
			content,
		})

		console.log({ respAdd: resp })

		const contents = await fetchContents(item.id)

		if (contents?.length === 0) {
			setContents([])
			setEditableContents(new Map())
			return
		}
		setEditableContents(() => new Map(contents?.map(content => [content.id, {
			id: content.id,
			originalValue: content.content,
			updatedValue: content.content,
			error: {
				isError: false,
				message: "",
			}
		}])))
		setContents(contents)
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

	const onDeleteContent = async (content) => {
		await deleteRequest(`http://localhost:8080/item/${selectedItem.id}/content/${content.id}/delete`).then(async (res) => {
			const contents = await fetchContents(selectedItem.id)
			setContents(contents)
		})
	}

	return (
		<div className={style.main}>
			<div className={style.groups}>
				<button style={{ margin: "20px 0px" }} className={style.addGroup} type="button">
					Add group
				</button>
				<div className={style.groupListItems}>
					{!error && groups?.map(group => (
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
						{contents?.map(content => (
							<div className={style.todoContainer} key={content.id}>
								<div className={style.contentContainer}>
									<input
										type="checkbox"
										checked={selectedContents.includes(content.id)}
										onChange={() => onCheckboxChange(content)}
									/>
									<textarea
										onChange={(ev) => onContentEdit(ev, content)}
										value={editableContents.get(content.id).updatedValue}
										className={style.contentInput}
									/>
								</div>
								<div>
									<DeleteIcon className={style.contentDelBtn} onClick={() => onDeleteContent(content)} />
								</div>
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