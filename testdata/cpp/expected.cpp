#include <iostream>
#include <vector>
#include <memory>

// TODO: Implement better error handling
// FIXME: Memory leak in constructor
// NOTE: This class is thread-safe

/*
   This is a multi-line comment
   that should be preserved
*/

/**
 * Documentation comment for the class
 * @param value Initial value
 */
class Example {
private:
    int value;

public:
    Example(int val) : value(val) {
        std::cout << "Creating Example with value: " << val << std::endl;
    }

    int getValue() const {
        return value;
    }

    void setValue(int newVal) {
        value = newVal;
    }

    /* Block comment in method */
    void processData() {
        std::vector<int> data = {1, 2, 3, 4, 5};

        for (const auto& item : data) {
            std::cout << item << " ";
        }
        std::cout << std::endl;
    }
};

int main() {
    Example ex(42);

    ex.processData();

    int val = ex.getValue();
    std::cout << "Value: " << val << std::endl;

    return 0;
}