export function HeaderComponent() {
    return (
        <div className="w-full flex items-center">
            <div className="border-2">
                ECOMMERCE
            </div>
            <div className="border-2">
                Delivery To Indonesia
            </div>
            <div className="flex items-center border-2">
                <div>select</div>
                <div><input type="text" /></div>
                <div><button type="button">search</button></div>
            </div>
            <div className="border-2">language</div>
            <div className="border-2">
                <div>hello, sign in</div>
                <div>Account & List</div>
            </div>
            <div className="border-2">
                <div>Return</div>
                <div>Order</div>
            </div>
            <div className="border-2">Cart</div>
        </div>
    )
}