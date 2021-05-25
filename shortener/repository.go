package shortener

type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	ListAll() (*[]Redirect, error)

	Store(redirect *Redirect) error
}
