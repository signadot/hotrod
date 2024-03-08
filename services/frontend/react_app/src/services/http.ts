

export const apiGet = async <R>(path: string) => {
    const d = await fetch(path)
    const json = await d.json()
    return json as R
}

export const apiPost = async <R>(path: string, data: Record<string, string> | object, headers: Record<string, string>) => {
    const d = await fetch(path, {
        headers: {
            ...headers,
            "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify(data),
    })
    const json = await d.json()
    return json as R
}

