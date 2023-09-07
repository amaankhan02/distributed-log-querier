import logging
import random
import time

# randomly pick between info, error, debug, and warning
 
# List of predefined log messages to randomly choose from
info_log_messages = [
    "Application started",
    "Data processing complete",
    "Service started successfully",
    "Configuration loaded",
    "File saved successfully",
    "Request received from client",
    "Cache cleared",
    "Database connection established",
    "Application version 2.0.1 is now running",
    "Connected to remote server at 192.168.1.100",
    "User 'alice' logged in from IP address 203.0.113.42",
    "Loaded configuration file 'config.yaml'",
    "Database initialized successfully",
    "Cache size: 256 MB",
    "Server started on port 8080",
    "Processing batch job: Importing data",
    "System uptime: 14 days, 6 hours",
    "Backup completed at 2:00 AM",
    "Log rotation completed successfully",
    "Service 'webapp' is now available",
    "User 'admin' updated user permissions",
    "API version 1.2.3 is now deprecated",
    "User 'guest' session timeout: 30 minutes",
    "Data export to 'export.csv' started"
]

debug_log_messages = [
    "Entering function: calculate_total",
    "Variable 'x' is set to 5",
    "Debugging information: ...",
    "Received API request: ...",
    "Debugging trace: ...",
    "Entering function 'calculate_total'",
    "Variable 'x' has value 42",
    "Debugging information for module 'payment'",
    "Received incoming request: POST /api/data",
    "Debugging trace for network layer",
    "Verbose debugging mode enabled",
    "Entering method 'process_request'",
    "HTTP request headers: {'Content-Type': 'application/json'}",
    "Querying database for user 'john_doe'",
    "Parsing JSON data: {'key': 'value'}",
    "Debugging output for component 'widget'",
    "Method 'validate_input' called with arguments: param1='value1', param2='value2'",
    "Socket connection established with host 'example.com'",
    "Debugging log entry #12345",
    "Traceback: File 'module.py', Line 42, Function 'method', Exception 'error message'",
    "Debugging message for development purposes",
    "Debugging data: {'debug_key': 'debug_value'}",
    "Entering 'try' block to handle exception",
    "Debugging step: Verifying data integrity"
]

error_log_messages = [
    "Critical error: database connection failed",
    "Authentication failed for user 'username'",
    "Unhandled exception: Division by zero",
    "Database query error: SQL syntax error",
    "Database query timeout",
    "Service unavailable",
    "Invalid input data: JSON parsing error",
    "Authentication failed for user 'admin'",
    "Unhandled exception: Division by zero",
    "Critical error: Unable to start the server",
    "Invalid input data: Missing required field 'email'",
    "SQL syntax error: Unexpected token near 'SELECT'",
    "File write error: Unable to save 'data.txt'",
    "HTTP 404 Not Found: Page '/page' does not exist",
    "Socket connection error: Connection refused",
    "Internal server error: Unable to process request",
    "Failed to open log file 'app.log'",
    "Resource allocation error: Memory exhausted",
    "Invalid request format: XML parsing error",
    "Network communication error: Connection reset",
    "Service unavailable: Maintenance in progress",
    "API request failed: HTTP 500 Internal Server Error",
    "Permission denied: User 'guest' cannot access 'admin' area",
    "Invalid configuration: Missing required settings",
    "Unhandled exception in module 'processing': Null reference"
]

warning_log_messages = [
    "Invalid input detected",
    "Low disk space warning",
    "Deprecated function in use",
    "Network connection unstable",
    "Resource usage exceeded warning",
    "Configuration option missing",
    "Low memory warning, available memory: 10 MB",
    "Deprecated API usage, consider updating",
    "Configuration file outdated, please update",
    "Resource usage near capacity, action required",
    "File deletion warning, 'temp.txt' will be removed",
    "Timeout in network communication, retrying",
    "Permission denied for action, access restricted",
    "Duplicate entry detected, data integrity at risk",
    "SSL certificate expiration warning, renew needed"

]

critical_log_messages = [
    "System crash: unrecoverable error",
    "Security breach detected",
    "Server outage: emergency shutdown",
    "Data corruption: critical issue",
    "Network failure: unable to connect",
    "Network infrastructure failure: disaster recovery plan activated"
    "Application terminated unexpectedly: core dump generated",
    "Server crash: Data loss imminent, urgent action required",
    "Security breach detected: Immediate response needed",
    "Power outage: Emergency shutdown initiated",
    "Data corruption detected: Data recovery in progress",
    "Network infrastructure failure: Disaster recovery activated",
    "Application terminated unexpectedly: Investigate immediately",
    "Database corruption: Data integrity compromised",
    "Hardware failure: Critical components malfunctioning",
    "Server overload: Service interruption imminent"
]

# list of predefined log levels to randomly choose from
log_levels = [logging.INFO, logging.DEBUG, logging.ERROR, logging.WARNING, logging.CRITICAL]

def generate_log_file(file, num_lines, log_with_time, known_lines, i):
    # Create a new logger instance with a unique name
    logger = logging.getLogger(str(i))

    # creates log file that logs DEBUG, INFO, WARNING, or ERROR (level = lowest possible level, setting it to debug means it logs all levels)
    # if file is already created it adds to the log file
    if log_with_time:
        logging.basicConfig(filename= file, level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s %(name)s')
    else:
        logging.basicConfig(filename= file, level=logging.DEBUG, format='%(levelname)s: %(message)s %(name)s')

    add_random_lines(num_lines, logger)
    add_known_lines(known_lines, logger)

def add_random_lines(num_lines, logger):
    log_message = ""

    for line_num in range(num_lines):

        # Choose a random log level
        log_level = random.choice(log_levels)

        if log_level == logging.INFO:
            log_message = random.choice(info_log_messages)
            logger.info(log_message)
        elif log_level == logging.DEBUG:
            log_message = random.choice(debug_log_messages)
            logger.debug(log_message)
        elif log_level == logging.ERROR:
            log_message = random.choice(error_log_messages)
            logger.error(log_message)
        elif log_level == logging.WARNING:
            log_message = random.choice(warning_log_messages)
            logger.warning(log_message)
        elif log_level == logging.CRITICAL:
            log_message = random.choice(critical_log_messages)
            logger.critical(log_message)

def add_known_lines(known_lines, logger):

    for (log_level, log_message) in known_lines:
        
        if log_level == logging.INFO:
            logger.info(log_message)
        elif log_level == logging.DEBUG:
            logger.debug(log_message)
        elif log_level == logging.ERROR:
            logger.error(log_message)
        elif log_level == logging.WARNING:
            logger.warning(log_message)
        elif log_level == logging.CRITICAL:
            logger.critical(log_message)

def main():

    files = ['test_log_file1.log', 'test_log_file2.log', 'test_log_file3.log']
    file_lines = [100, 56, 456]
    known_lines = [[(logging.info, "User logged in"),
                    (logging.error, "Server responded with HTTP error code: 404"), 
                    (logging.warning, "Disk space running low, clear some space"),
                    (logging.critical, "Critical system failure: Application halted")],

                   [(logging.info, "Email sent to 'user@example.com'",), 
                    (logging.info,  "Disk space usage: 80 percent used")],

                   [(logging.error, "File not found: 'file.txt'"),
                    (logging.debug, "Debug message: Processing step 3"),
                    (logging.info, "Scheduled maintenance task started"),
                    (logging.info, "Received 200 OK response from API endpoint"),
                    (logging.critical, "Application halted: fatal error")]]
    
    with_time = [False, True, True]


    for i in range(len(files)):
        generate_log_file(files[i], file_lines[i], with_time[i], known_lines[i], i)
    
    logging.shutdown()


if __name__ == "__main__":
    main()









