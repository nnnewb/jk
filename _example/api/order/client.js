var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) {
        return value instanceof P ? value : new P(function (resolve) {
            resolve(value);
        });
    }

    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) {
            try {
                step(generator.next(value));
            } catch (e) {
                reject(e);
            }
        }

        function rejected(value) {
            try {
                step(generator["throw"](value));
            } catch (e) {
                reject(e);
            }
        }

        function step(result) {
            result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected);
        }

        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export default {
    baseURL: "",
    cancel_order: function (payload, init) {
        return __awaiter(this, void 0, void 0, function* () {
            const u = new URL("/api/v1/order-service/order/cancel", this.baseURL);
            init.body = JSON.stringify(payload);
            init.method = "POST";
            const req = new Request(u, init);
            const resp = yield fetch(req, init);
            return yield resp.json();
        });
    },
    create_order: function (payload, init) {
        return __awaiter(this, void 0, void 0, function* () {
            const u = new URL("/api/v1/order-service/order", this.baseURL);
            init.body = JSON.stringify(payload);
            init.method = "POST";
            const req = new Request(u, init);
            const resp = yield fetch(req, init);
            return yield resp.json();
        });
    },
    order_detail: function (payload, init) {
        return __awaiter(this, void 0, void 0, function* () {
            const u = new URL("/api/v1/order-service/order/detail", this.baseURL);
            Object.getOwnPropertyNames(payload).map((prop) => u.searchParams.append(prop, payload[prop]));
            init.method = "GET";
            const req = new Request(u, init);
            const resp = yield fetch(req, init);
            return yield resp.json();
        });
    },
    update: function (payload, init) {
        return __awaiter(this, void 0, void 0, function* () {
            const u = new URL("/api/v1/order-service/order", this.baseURL);
            init.body = JSON.stringify(payload);
            init.method = "PUT";
            const req = new Request(u, init);
            const resp = yield fetch(req, init);
            return yield resp.json();
        });
    },
};
//# sourceMappingURL=client.js.map