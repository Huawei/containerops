import unittest
import xmlrunner

unittest.main(
        module=None,
        testRunner=xmlrunner.XMLTestRunner(output="/tmp/output"),
        failfast=False, buffer=False, catchbreak=False)
