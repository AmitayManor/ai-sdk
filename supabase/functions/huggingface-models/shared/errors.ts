export class ModelError extends Error {
    constructor(
        message: string,
        public statusCode: number = 500
    ) {
        super(message)
        this.name = 'ModelError'
    }
}

export class ValidationError extends ModelError {
    constructor(message: string) {
        super(message, 400)
        this.name = 'ValidationError'
    }
}

export class ProcessingError extends ModelError {
    constructor(message: string) {
        super(message, 500)
        this.name = 'ProcessingError'
    }
}

export class TimeoutError extends ModelError {
    constructor(message: string = 'Request timed out') {
        super(message, 408)
        this.name = 'TimeoutError'
    }
}

export function handleError(error: unknown): Response {
    const modelError = error instanceof ModelError
        ? error
        : new ProcessingError(error instanceof Error ? error.message : 'Unknown error')

    return new Response(
        JSON.stringify({
            error: modelError.message,
            status: 'error'
        }),
        {
            status: modelError.statusCode,
            headers: { 'Content-Type': 'application/json' }
        }
    )
}