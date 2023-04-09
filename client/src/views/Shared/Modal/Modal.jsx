import React from "react"
import { Overlay } from "../Overlay"

import styles from "./styles.module.css"


const Modal = ({
	// eslint-disable-next-line react/prop-types
	children
}) => {
	return (
		<Overlay>
			<div className={styles.container}>
				<div>
					{children}
				</div>
			</div>
		</Overlay>
	)
}

export default Modal