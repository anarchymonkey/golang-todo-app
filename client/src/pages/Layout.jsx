import React, { memo } from "react"
import { Header } from "./"
import Router from "./Router"


const Layout = () => {

	return (
		<div>
			<Header />
			<Router />
		</div>
	)
}

export default memo(Layout)