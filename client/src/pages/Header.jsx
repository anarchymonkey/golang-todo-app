import React from "react"
import styles from "./header.module.css"

function Header() {

	const goToPage = (pathName) => {
		console.log("Going to page", pathName)
	}

	return (
		<div className={styles.main}>
			<div className={styles.left}>
				Daily Shenanigans
			</div>
			<div className={styles.right}>
				<div onClick={() => goToPage("/collections")}>Collections</div>
				<div onClick={() => goToPage("/about")}>About</div>
			</div>
		</div>
	)
}

export default Header