


const currencyFormatter = new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR'
})

export function formatCents(cents: number): string {
    return currencyFormatter.format(cents / 100)
}

export function parseEurosToCents(input: string): number | null {
    const normalized = input.trim().replace(',', '.')
    if (normalized === '') {
        return null
    }

    const euros = Number(normalized)
    if (!Number.isFinite(euros) || euros <= 0) {
        return null
    }

    return Math.round(euros * 100)
}