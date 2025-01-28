"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var _a, _b, _c, _d, _e;
Object.defineProperty(exports, "__esModule", { value: true });
// import { serve } from '$std/http/server.ts'
// import { createClient } from '@supabase/supabase-js'
var huggingface_ts_1 = require("./../../shared/huggingface.ts");
var errors_ts_1 = require("./../../shared/errors.ts");
require("jsr:@supabase/functions-js/edge-runtime.d.ts");
// Get environment variables
var env = {
    SUPABASE_URL: (_a = Deno.env.get('SUPABASE_URL')) !== null && _a !== void 0 ? _a : '',
    SUPABASE_ANON_KEY: (_b = Deno.env.get('SUPABASE_ANON_KEY')) !== null && _b !== void 0 ? _b : '',
    HF_API_TOKEN: (_c = Deno.env.get('HF_API_TOKEN')) !== null && _c !== void 0 ? _c : '',
    MODEL_TIMEOUT: Number((_d = Deno.env.get('MODEL_TIMEOUT')) !== null && _d !== void 0 ? _d : 30000),
    MAX_RETRIES: Number((_e = Deno.env.get('MAX_RETRIES')) !== null && _e !== void 0 ? _e : 3)
};
serve(function (req) { return __awaiter(void 0, void 0, void 0, function () {
    var _a, id, input, modelId, hfClient, supabase, startTime, result, processingTime, response, error_1;
    return __generator(this, function (_b) {
        switch (_b.label) {
            case 0:
                _b.trys.push([0, 6, , 7]);
                return [4 /*yield*/, req.json()];
            case 1:
                _a = _b.sent(), id = _a.id, input = _a.input, modelId = _a.modelId;
                hfClient = new huggingface_ts_1.HuggingFaceClient(env);
                supabase = createClient(env.SUPABASE_URL, env.SUPABASE_ANON_KEY);
                // Validate input
                if (!input || typeof input !== 'string') {
                    throw new Error('Invalid input: must be a non-empty string');
                }
                // Validate model
                return [4 /*yield*/, hfClient.validateModel(modelId)
                    // Update request status
                ];
            case 2:
                // Validate model
                _b.sent();
                // Update request status
                return [4 /*yield*/, supabase
                        .from('model_requests')
                        .update({ status: 'IN_PROGRESS' })
                        .eq('id', id)
                    // Process request
                ];
            case 3:
                // Update request status
                _b.sent();
                startTime = Date.now();
                return [4 /*yield*/, hfClient.textToText(modelId, input)];
            case 4:
                result = _b.sent();
                processingTime = Date.now() - startTime;
                response = {
                    output: result,
                    processingTime: processingTime,
                    status: 'success',
                    tokenCount: result.split(/\s+/).length // Basic token count
                };
                // Update request with results
                return [4 /*yield*/, supabase
                        .from('model_requests')
                        .update({
                        status: 'COMPLETED',
                        output_data: response,
                        completed_at: new Date().toISOString(),
                        processing_time: processingTime,
                        token_count: response.tokenCount
                    })
                        .eq('id', id)];
            case 5:
                // Update request with results
                _b.sent();
                return [2 /*return*/, new Response(JSON.stringify(response), {
                        headers: { 'Content-Type': 'application/json' }
                    })];
            case 6:
                error_1 = _b.sent();
                return [2 /*return*/, (0, errors_ts_1.handleError)(error_1)];
            case 7: return [2 /*return*/];
        }
    });
}); });
