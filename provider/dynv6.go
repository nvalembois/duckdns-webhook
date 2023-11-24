// https://dynv6.com/users/sign_up

package provider

// func main() {
// 	// Configuration de l'authentification
// 	config := &ssh.ClientConfig{
// 		User: "your-username",
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password("your-password"),
// 			// Vous pouvez également utiliser d'autres méthodes d'authentification, comme la clé privée, en fonction de votre configuration.
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Ignorer la vérification de l'empreinte du serveur (attention en production !)
// 	}

// 	// Connexion au serveur SSH
// 	client, err := ssh.Dial("tcp", "your-server-address:22", config)
// 	if err != nil {
// 		fmt.Println("Erreur lors de la connexion au serveur SSH :", err)
// 		os.Exit(1) // Code de retour non nul pour indiquer une erreur
// 	}
// 	defer client.Close()

// 	// Exécution de la commande à distance
// 	session, err := client.NewSession()
// 	if err != nil {
// 		fmt.Println("Erreur lors de la création de la session SSH :", err)
// 		os.Exit(1) // Code de retour non nul pour indiquer une erreur
// 	}
// 	defer session.Close()

// 	// Vous pouvez modifier la commande que vous souhaitez exécuter
// 	command := "ls -l"

// 	// Redirection des entrées/sorties de la session SSH vers les entrées/sorties de votre programme
// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr
// 	session.Stdin = os.Stdin

// 	// Exécution de la commande à distance
// 	err = session.Run(command)
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'exécution de la commande :", err)
// 		os.Exit(1) // Code de retour non nul pour indiquer une erreur
// 	}

// 	// La fonction main renvoie 0 pour indiquer une exécution réussie
// 	os.Exit(0)
// }
