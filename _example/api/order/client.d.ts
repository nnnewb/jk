interface CancelOrderRequest {
    order_id: string;
}

interface CancelOrderResponse {
    code: number;
    message: string;
}

interface OrderItem {
    item_id: string;
    quantity: number;
}

interface CreateOrderRequest {
    order_info: Array<OrderItem>;
}

interface CreateOrderResponse {
    code: number;
    message: string;
    order_id: string;
}

interface GetOrderDetailRequest {
    order_id: string;
}

interface GetOrderDetailResponse {
    code: number;
    message: string;
    order_info: Array<OrderItem>;
}

interface UpdateOrderRequest {
    order_info: Array<OrderItem>;
}

interface UpdateOrderResponse {
    code: number;
    message: string;
}

declare const _default: {
    baseURL: string;
    cancel_order: (payload: CancelOrderRequest, init?: RequestInit) => Promise<CancelOrderResponse>;
    create_order: (payload: CreateOrderRequest, init?: RequestInit) => Promise<CreateOrderResponse>;
    order_detail: (payload: GetOrderDetailRequest, init?: RequestInit) => Promise<GetOrderDetailResponse>;
    update: (payload: UpdateOrderRequest, init?: RequestInit) => Promise<UpdateOrderResponse>;
};
export default _default;
//# sourceMappingURL=client.d.ts.map