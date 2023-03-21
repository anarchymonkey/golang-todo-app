import { memo } from 'react'


const Layout = ({
    children
}) => {

    return (
        <div>
            {children}
        </div>
    )
}

export default memo(Layout);