


const currencyFormatter = new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR'
})

export function formatCents(cents: number): string {
    return currencyFormatter.format(cents / 100)
}