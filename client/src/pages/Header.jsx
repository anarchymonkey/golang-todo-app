

const Header = () => {

    const goToPage = (pathName) => {
        console.log("Going to page", pathName)
    }

    return (
        <div className="main">
            <div className="left">
                Daily Shenanigans
            </div>
            <div className="right">
                <div onClick={() => goToPage('/collections')}>Collections</div>
                <div onClick={() => goToPage('/about')}>About</div>
            </div>
        </div>
    )
}

export default Header;