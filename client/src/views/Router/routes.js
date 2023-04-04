import { lazy } from "react"

// constants
import { PATHS } from "../../constants"
// import { Todo } from "../../pages"



const Collections = lazy(() => import("../../pages/Todo.jsx"))

console.log("Collections", Collections)

export const getRoutes = () => {
	return [{
		path: PATHS.HOME,
		component: Collections,
	},{
		path: PATHS.COLLECTIONS,
		component: Collections,
	}]
}