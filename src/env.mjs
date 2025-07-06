import { spawnSync } from "child_process";
export class EnvError extends Error {
    constructor(code, message) {
        super(message);
        this.name = 'EnvError';
        this.code = code;
        this.message = message;
    }
}
export function execPwshCommand(command) {
    var _a;
    const result = spawnSync('powershell', ['-c', command]);
    if (result.status !== 0) {
        throw new EnvError((_a = result.status) !== null && _a !== void 0 ? _a : 1, result.stderr.toString().trim());
    }
    return result.stdout.toString().trim();
}
export function readEnv(name, scope) {
    return execPwshCommand(`[System.Environment]::GetEnvironmentVariable("${name}", [System.EnvironmentVariableTarget]::${scope})`);
}
export function writeEnv(name, value, scope) {
    return execPwshCommand(`[System.Environment]::SetEnvironmentVariable("${name}", "${value}", [System.EnvironmentVariableTarget]::${scope})`);
}
export function addToPath(scope, ...paths) {
    const existingPath = readEnv('PATH', scope);
    const existingPathArray = existingPath.split(';').filter((p) => p);
    const newPath = [...new Set([...existingPathArray, ...paths])].join(';');
    return writeEnv('PATH', newPath, scope);
}
