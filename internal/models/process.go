package models

type UploadedFile struct {
	Filename    string             // Le nom du fichier
	ContentType string             // Le type MIME du fichier (ex. "application/pdf")
	Save        func(string) error // Fonction pour sauvegarder le fichier à l'emplacement donné
}

type UpdateRequest struct {
	Name string   // Le nouveau nom du document
	Type string   // Le nouveau type MIME du document
	Tags []string // Les nouveaux tags du document
}
