// Code generated by jk generate all -t OrderService; DO NOT EDIT.
export default {
    base_url: "",
    /**
     * send request
     * @param method
     * @param url
     * @param config {{headers?: HeadersInit, params?: URLSearchParams, data?: any}} request configuration
     * @returns { Promise<Response> }
     */
    ajax: function (method, url, config) {
        method = method.toLowerCase();
        let u = null;
        if (this.base_url && this.base_url.length !== 0) {
            u = new URL(url, this.base_url)
        } else {
            u = new URL(url);
        }

        let params = undefined;
        if (config.params) {
            params = config.params;
        }

        u.searchParams = params;

        let body = undefined;
        if (config.data) {
            body = JSON.stringify(config.data);
        }

        let headers = undefined;
        if (config.headers) {
            headers = config.headers;
        }

        return fetch(u, {method, headers, body});
    },
    /**
     *
     * @param payload {{order_id:string,}}
     * @param config {{headers?: HeadersInit}} request configuration
     */
    cancel_order: function (payload, config) {
        return this.ajax(
            'POST',
            '/api/v1/order-service/order/cancel',
            config || {},
        );
    },
    /**
     *
     * @param payload {{order_info:Array<{item_id:string,quantity:number,}>,}}
     * @param config {{headers?: HeadersInit}} request configuration
     */
    create_order: function (payload, config) {
        return this.ajax(
            'POST',
            '/api/v1/order-service/order',
            config || {},
        );
    },
    /**
     *
     * @param payload {{order_id:string,}}
     * @param config {{headers?: HeadersInit}} request configuration
     */
    order_detail: function (payload, config) {
        return this.ajax(
            'GET',
            '/api/v1/order-service/order/detail',
            config || {},
        );
    },
    /**
     *
     * @param payload {{order_info:Array<{item_id:string,quantity:number,}>,}}
     * @param config {{headers?: HeadersInit}} request configuration
     */
    update: function (payload, config) {
        return this.ajax(
            'PUT',
            '/api/v1/order-service/order',
            config || {},
        );
    },
};