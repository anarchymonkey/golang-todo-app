
import React, { Suspense } from "react"
import { Switch } from "react-router-dom"



const Router = () => {

	return (
		<Suspense fallback={<div>Loading</div>}>
			<Switch>

			</Switch>
		</Suspense>
	)
}

export default Router