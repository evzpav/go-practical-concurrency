# Go Practical Concurrency Examples

Compilation of practical and simple examples of concurrency in Golang.

- a) Processing Text Files (e.g.:.csv, .xslx, .txt):
    1. Process text file
        * Read whole text file
        * Process all lines (do something with that text) (conc)
        * Write result to output file in the same order
        ```
            cd a.text_files/1.process_file
            go run main.go
        ```
    2. Process text file line by line - unordered results
        * Read text file line by line (conc)
        * Process each line (do something) (conc)
        * Write line by line (out of order) (conc)
        ```
            cd a.text_files/2.read_process_write_line_by_line
            go run main.go
        ```
    3. Process text file - ordered results
        * Read text file line by line (conc)
        * Process each line (do something) (conc)
        * Order results (conc)
        * Write line by line (conc)
        ```
            cd a.text_files/3.read_process_write_line_by_line_ordered
            go run main.go
        ```
- b) Processing Multiple HTTP requests:
    1. Site status checker
        * Iterate through array of urls
        * Make request to check if error or status code of the response (conc)
        * Print result to console
        ```
           cd b.http_requests/1.site_status_checker
           go run main.go
        ```

    2. Download big file
        * Make request to get "Content-Length" of file
        * Split size in sessions of your choice
        * Make multiple requests to retrieve each session(conc)
        * Save each session to temp files
        * Merge temp files into one file

        ```
           cd b.http_requests/2.download_file
           go run main.go
        ```
    3.  Make multiple requests get data and save to file
        * Make request to get content (conc)
        * Make request for each content (conc)
        * Save to text file (conc)

        ```
            cd b.http_requests/3.get_hacker_news_stories
            go run main.go
        ```

PS: Probably there are many other ways (better) to do the same things. Appreciate suggestions of improvements or new examples to be added.