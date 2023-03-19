
import "./styles.css"

const Card = ({
    children
}) => {
    return (
        <div className="container">
            {children}
        </div>
    )
}

export default Card;