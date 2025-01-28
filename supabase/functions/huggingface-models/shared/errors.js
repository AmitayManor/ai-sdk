"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.handleError = exports.TimeoutError = exports.ProcessingError = exports.ValidationError = exports.ModelError = void 0;
var ModelError = /** @class */ (function (_super) {
    __extends(ModelError, _super);
    function ModelError(message, statusCode) {
        if (statusCode === void 0) { statusCode = 500; }
        var _this = _super.call(this, message) || this;
        _this.statusCode = statusCode;
        _this.name = 'ModelError';
        return _this;
    }
    return ModelError;
}(Error));
exports.ModelError = ModelError;
var ValidationError = /** @class */ (function (_super) {
    __extends(ValidationError, _super);
    function ValidationError(message) {
        var _this = _super.call(this, message, 400) || this;
        _this.name = 'ValidationError';
        return _this;
    }
    return ValidationError;
}(ModelError));
exports.ValidationError = ValidationError;
var ProcessingError = /** @class */ (function (_super) {
    __extends(ProcessingError, _super);
    function ProcessingError(message) {
        var _this = _super.call(this, message, 500) || this;
        _this.name = 'ProcessingError';
        return _this;
    }
    return ProcessingError;
}(ModelError));
exports.ProcessingError = ProcessingError;
var TimeoutError = /** @class */ (function (_super) {
    __extends(TimeoutError, _super);
    function TimeoutError(message) {
        if (message === void 0) { message = 'Request timed out'; }
        var _this = _super.call(this, message, 408) || this;
        _this.name = 'TimeoutError';
        return _this;
    }
    return TimeoutError;
}(ModelError));
exports.TimeoutError = TimeoutError;
function handleError(error) {
    var modelError = error instanceof ModelError
        ? error
        : new ProcessingError(error instanceof Error ? error.message : 'Unknown error');
    return new Response(JSON.stringify({
        error: modelError.message,
        status: 'error'
    }), {
        status: modelError.statusCode,
        headers: { 'Content-Type': 'application/json' }
    });
}
exports.handleError = handleError;
