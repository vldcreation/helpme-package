class SomeError extends Error {
  constructor(message) {
    super(message);
    this.name = 'SomeError';
  }
}

// color constants
export const RESET = '\x1b[0m';
export const RED = '\x1b[31m';
const GREEN = '\x1b[32m';
const YELLOW = '\x1b[33m';

export default SomeError;