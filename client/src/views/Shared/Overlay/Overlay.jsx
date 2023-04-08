/* eslint-disable react/prop-types */
/* eslint-disable indent */
import React from "react"
import styles from "./styles.module.css"


const Overlay = ({
    children,
}) => {

	return (
		<div className={styles.overlay}>
			{children}
		</div>
	)
}

export default Overlay