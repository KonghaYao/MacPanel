const MCP_SERVER_NAME = 'macpanel';

export const MCP_API_KEY_PLACEHOLDER = '<your-api-key>';

export function getMcpEndpointUrl(): string {
    return `${window.location.origin}/mcp`;
}

function buildMcpConfig(secret: string) {
    return {
        mcpServers: {
            [MCP_SERVER_NAME]: {
                url: getMcpEndpointUrl(),
                headers: {
                    Authorization: `Bearer ${secret}`,
                },
            },
        },
    };
}

export function buildMcpCursorConfig(secret: string): string {
    return JSON.stringify(buildMcpConfig(secret), null, 2);
}

export function buildMcpExampleConfig(): string {
    return buildMcpCursorConfig(MCP_API_KEY_PLACEHOLDER);
}
