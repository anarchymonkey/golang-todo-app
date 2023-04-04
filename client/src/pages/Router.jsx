
import React, { Suspense } from "react"
import { Route, Routes } from "react-router-dom"
import { getRoutes } from "../views/Router/routes"

const Router = () => {
	// console.log(getRoutes())
	return (
		<Routes>
			{getRoutes().map(route => (
				<Route key={route.path} path={route.path} element={
					<Suspense fallback={<div>Loading</div>}>
						<route.component />
					</Suspense>
				} />
			))}
		</Routes>
	)
}

export default Router
