import { memo } from 'react'


const Layout = ({
    children
}) => {

    return (
        <div>
            <h3>This is the layout module</h3>
            {children}
        </div>
    )
}

export default memo(Layout);