package shortener

type RedirectService interface {
	Find(code string) (*Redirect, error)
	ListAll() (*[]Redirect, error)

	Store(redirect *Redirect) error
}
