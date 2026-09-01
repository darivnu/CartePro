import { useLogout } from '../../auth/useAuth'

export function PartnerHome() {
  const logout = useLogout()

  return (
    <div className="p-4">
      Partner area
      <button onClick={() => logout.mutate()}>Log out</button>
    </div>
  )
}
